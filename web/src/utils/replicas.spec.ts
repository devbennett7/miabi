import { describe, expect, it } from 'vitest'
import { serviceReplicas } from './replicas'

describe('serviceReplicas', () => {
  it('is null for a container app', () => {
    expect(serviceReplicas({ runtime_kind: 'container', replicas: 1 })).toBeNull()
  })

  it('reports running over desired, degraded when short', () => {
    expect(serviceReplicas({ runtime_kind: 'service', replicas: 3, running_replicas: 2, desired_replicas: 3 }))
      .toEqual({ running: 2, desired: 3, degraded: true })
    expect(serviceReplicas({ runtime_kind: 'service', replicas: 2, running_replicas: 2, desired_replicas: 2 }))
      .toEqual({ running: 2, desired: 2, degraded: false })
  })

  it('falls back to the configured count when the swarm could not be asked', () => {
    expect(serviceReplicas({ runtime_kind: 'service', replicas: 2 })).toEqual({ running: null, desired: 2, degraded: false })
  })
})
