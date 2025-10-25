import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Trello CLI',
  description: 'A comprehensive Trello CLI tool optimized for LLM integration',
  
  head: [
    ['link', { rel: 'icon', href: '/logo.png' }]
  ],

  themeConfig: {
    logo: '/logo.png',
    
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/installation' },
      { text: 'Reference', link: '/reference/commands' },
      { text: 'Examples', link: '/examples/api-integration' }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Authentication', link: '/guide/authentication' },
            { text: 'Quick Start', link: '/guide/quick-start' }
          ]
        }
      ],
      '/reference/': [
        {
          text: 'Command Reference',
          items: [
            { text: 'Commands Overview', link: '/reference/commands' },
            { text: 'Global Flags', link: '/reference/flags' },
            { text: 'Boards', link: '/reference/boards' },
            { text: 'Lists', link: '/reference/lists' },
            { text: 'Cards', link: '/reference/cards' },
            { text: 'Labels', link: '/reference/labels' },
            { text: 'Checklists', link: '/reference/checklists' },
            { text: 'Members', link: '/reference/members' },
            { text: 'Attachments', link: '/reference/attachments' },
            { text: 'Batch Operations', link: '/reference/batch' },
            { text: 'Configuration', link: '/reference/config' }
          ]
        }
      ],
      '/examples/': [
        {
          text: 'Examples & Tutorials',
          items: [
            { text: 'API Integration', link: '/examples/api-integration' },
            { text: 'LLM Workflows', link: '/examples/llm-workflows' },
            { text: 'Automation', link: '/examples/automation' },
            { text: 'Use Cases', link: '/examples/use-cases' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/danbruder/trello-cli' }
    ],

    search: {
      provider: 'local'
    }
  }
})
