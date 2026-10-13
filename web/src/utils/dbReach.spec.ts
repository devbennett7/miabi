import { describe, expect, it } from 'vitest'
import { canReachDatabase } from './dbReach'

describe('canReachDatabase', () => {
  const overlayDB = { server_id: 2, cluster_id: 0, network_name: 'mb-ws2-ec333153e7b0-ov' }
  const bridgeDB = { server_id: 2, cluster_id: 0, network_name: 'mb-ws2-ec333153e7b0' }

  it('lets any node of the cluster reach a database on a workspace overlay', () => {
    expect(canReachDatabase({ server_id: 0, cluster_id: 0 }, overlayDB)).toBe(true)
    expect(canReachDatabase({ server_id: 5 }, overlayDB)).toBe(true)
  })

  it('keeps an overlay database to its own cluster', () => {
    expect(canReachDatabase({ server_id: 2, cluster_id: 3 }, overlayDB)).toBe(false)
  })

  it('keeps a bridge database to its own node', () => {
    expect(canReachDatabase({ server_id: 2 }, bridgeDB)).toBe(true)
    expect(canReachDatabase({ server_id: 0 }, bridgeDB)).toBe(false)
  })

  it('treats a database on no workspace network as node-local', () => {
    expect(canReachDatabase({ server_id: 1 }, { server_id: 0 })).toBe(false)
    expect(canReachDatabase({}, {})).toBe(true)
  })
})
