import { defineWorkspace } from 'vitest/config'

export default defineWorkspace([
  {
    test: {
      name: 'ui',
      root: './packages/ui',
      environment: 'jsdom',
      setupFiles: ['../../test-setup.ts'],
    },
  },
  {
    test: {
      name: 'api-client',
      root: './packages/api-client',
      environment: 'jsdom',
      setupFiles: ['../../test-setup.ts'],
    },
  },
  {
    test: {
      name: 'web-professional',
      root: './apps/web-professional',
      environment: 'jsdom',
      setupFiles: ['../../test-setup.ts'],
    },
  },
])
