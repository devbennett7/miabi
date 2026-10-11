<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { authApi } from '@/api/auth'
import { workspaceApi } from '@/api/workspaces'
import { apiErrorMessage } from '@/api/client'
import AuthShell from './AuthShell.vue'
import type { Brand, InvitationPreview } from '@/api/types'
import { useAuthStore } from '@/stores/auth'
import { useWorkspaceStore } from '@/stores/workspace'
import { useThemeStore } from '@/stores/theme'
import { useBrandStore } from '@/stores/brand'
import { savePendingInvite, pendingInvite, clearPendingInvite } from '@/utils/pendingInvite'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()
const ws = useWorkspaceStore()
const theme = useThemeStore()
const brandStore = useBrandStore()

const brand = ref<Brand>({})
const brandName = computed(() => brand.value.name?.trim() || 'Miabi')
const brandLogo = computed(() => (theme.isDark && brand.value.logo_dark_url?.trim()) || brand.value.logo_url?.trim() || '/brand/miabi-mark.svg')

const token = ref('')
const preview = ref<InvitationPreview | null>(null)
const loading = ref(true)
const invalid = ref(false)
const busy = ref(false)
const error = ref('')
const name = ref('')
const password = ref('')

// Matches handlers.MinRegistrationPasswordLen.
const MIN_PASSWORD = 12

const signedInAsInvitee = computed(
  () => !!auth.user && !!preview.value && auth.user.email.toLowerCase() === preview.value.email.toLowerCase(),
)

const title = computed(() => {
  if (invalid.value) return t('acceptInvite.invalidTitle')
  if (!preview.value) return t('acceptInvite.loading')
  return t('acceptInvite.title', { workspace: preview.value.workspace_name })
})

const subtitle = computed(() => {
  const p = preview.value
  if (!p) return ''
  return p.invited_by_name
    ? t('acceptInvite.subtitleBy', { name: p.invited_by_name, role: p.role })
    : t('acceptInvite.subtitle', { role: p.role })
})

onMounted(async () => {
  try {
    const { data } = await authApi.status()
    brand.value = data.data?.brand ?? {}
    brandStore.set(brand.value)
  } catch {
    // Branding is cosmetic here.
  }
  theme.applyBrandAccent(brand.value.accent)

  const fromLink = typeof route.query.token === 'string' ? route.query.token.trim() : ''
  token.value = fromLink || pendingInvite() || ''
  if (fromLink) {
    savePendingInvite(fromLink)
    // Keep the token out of the address bar, history and any Referer once it is held.
    router.replace({ name: 'accept-invitation' })
  }
  if (!token.value) {
    invalid.value = true
    loading.value = false
    return
  }
  try {
    const { data } = await authApi.previewInvitation(token.value)
    preview.value = data.data
  } catch {
    invalid.value = true
    clearPendingInvite()
  } finally {
    loading.value = false
  }
})

async function enterWorkspace(id: number) {
  clearPendingInvite()
  await ws.fetchWorkspaces()
  ws.setWorkspace(id)
  router.replace('/')
}

async function accept() {
  busy.value = true
  error.value = ''
  try {
    const { data } = await workspaceApi.acceptInvitationToken(token.value)
    await enterWorkspace(data.data.id)
  } catch (e) {
    error.value = apiErrorMessage(e, t('acceptInvite.acceptFailed'))
  } finally {
    busy.value = false
  }
}

function signIn() {
  router.push({ name: 'login' })
}

async function signOut() {
  await auth.logout()
}

async function signUp() {
  error.value = ''
  if (password.value.length < MIN_PASSWORD) {
    error.value = t('acceptInvite.passwordTooShort', { n: MIN_PASSWORD })
    return
  }
  const p = preview.value
  if (!p) return
  busy.value = true
  try {
    const { data } = await authApi.registerFromInvitation({
      token: token.value, name: name.value.trim(), password: password.value,
    })
    const result = await auth.login(p.email, password.value)
    if (result !== 'ok') {
      clearPendingInvite()
      router.replace({ name: 'login' })
      return
    }
    await enterWorkspace(data.data.workspace_id)
  } catch (e) {
    error.value = apiErrorMessage(e, t('acceptInvite.signupFailed'))
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <AuthShell
    hero
    :brand-name="brand.name ? brandName : undefined"
    :brand-logo="brandLogo"
    :title="title"
    :subtitle="invalid ? '' : subtitle"
    :error="error"
  >
    <p v-if="loading" class="invite-text"><span class="mdi mdi-loading mdi-spin"></span></p>

    <template v-else-if="invalid">
      <p class="invite-text">{{ $t('acceptInvite.invalidBody') }}</p>
      <RouterLink :to="auth.isAuthenticated ? '/' : { name: 'login' }" class="btn btn-primary invite-submit">
        {{ auth.isAuthenticated ? $t('acceptInvite.openConsole') : $t('register.goToSignIn') }}
      </RouterLink>
    </template>

    <template v-else-if="preview && auth.isAuthenticated">
      <template v-if="signedInAsInvitee">
        <button class="btn btn-primary invite-submit" :disabled="busy" @click="accept">
          <span v-if="busy" class="mdi mdi-loading mdi-spin"></span>
          {{ $t('acceptInvite.join', { workspace: preview.workspace_name }) }}
        </button>
      </template>
      <template v-else>
        <p class="invite-text">
          <i18n-t keypath="acceptInvite.otherAccount" tag="span">
            <template #invited><strong>{{ preview.email }}</strong></template>
            <template #current><strong>{{ auth.user?.email }}</strong></template>
          </i18n-t>
        </p>
        <button class="btn btn-primary invite-submit" @click="signOut">{{ $t('acceptInvite.signOut') }}</button>
      </template>
    </template>

    <template v-else-if="preview && preview.account_exists">
      <p class="invite-text">
        <i18n-t keypath="acceptInvite.signInToAccept" tag="span"><template #email><strong>{{ preview.email }}</strong></template></i18n-t>
      </p>
      <button class="btn btn-primary invite-submit" @click="signIn">{{ $t('acceptInvite.signIn') }}</button>
    </template>

    <template v-else-if="preview && preview.sso_required">
      <p class="invite-text">{{ $t('acceptInvite.ssoRequired') }}</p>
      <button class="btn btn-primary invite-submit" @click="signIn">{{ $t('acceptInvite.continueWithSso') }}</button>
    </template>

    <form v-else-if="preview && preview.signup_allowed" class="auth-form" @submit.prevent="signUp">
      <div class="form-group">
        <label class="form-label" for="inv-email">{{ $t('register.email') }}</label>
        <input id="inv-email" :value="preview.email" type="email" class="form-input" readonly />
      </div>
      <div class="form-group">
        <label class="form-label" for="inv-name">{{ $t('apps.form.name') }}</label>
        <input id="inv-name" v-model="name" class="form-input" autocomplete="name" :disabled="busy" required />
      </div>
      <div class="form-group">
        <label class="form-label" for="inv-password">{{ $t('register.password') }}</label>
        <input id="inv-password" v-model="password" type="password" class="form-input" autocomplete="new-password"
          :minlength="MIN_PASSWORD" :disabled="busy" required />
        <p class="form-hint">{{ $t('acceptInvite.passwordHint', { n: MIN_PASSWORD }) }}</p>
      </div>
      <button class="btn btn-primary invite-submit" :disabled="busy">
        <span v-if="busy" class="mdi mdi-loading mdi-spin"></span>
        {{ $t('acceptInvite.createAndJoin') }}
      </button>
    </form>

    <p v-else class="invite-text">
      <i18n-t keypath="acceptInvite.noSignup" tag="span"><template #email><strong>{{ preview?.email }}</strong></template></i18n-t>
    </p>

    <template #footer>
      <RouterLink v-if="!auth.isAuthenticated" :to="{ name: 'login' }" class="invite-back">
        <span class="mdi mdi-arrow-left"></span>{{ $t('login.backToSignIn') }}</RouterLink>
    </template>
  </AuthShell>
</template>

<style scoped>
.invite-submit {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.invite-text {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 16px;
}
.form-hint { font-size: 12px; color: var(--text-muted); margin-top: 4px; }
.invite-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-muted);
  text-decoration: none;
}
.invite-back:hover { color: var(--primary-600); }
</style>
