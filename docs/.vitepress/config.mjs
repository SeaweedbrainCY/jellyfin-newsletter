import { defineConfig } from 'vitepress'

export default defineConfig({
    title: 'Jellyfin Newsletter',
    description: 'Documentation for Jellyfin Newsletter',
    lang: 'en-US',
    lastUpdated: true,
    cleanUrls: true,
    head: [['link', { rel: 'icon', href: '/favicon.png' }]],
    themeConfig: {
        logo: "https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/jellyfin_newsletter.png",
        nav: [
            { text: 'Documentation', link: '/docs' },
        ],
        search: {
            provider: 'local'
        },
        footer: {
            message: 'Released under the AGPL-3.0 license License.',
            copyright: 'Copyright 2025-present Nathan Stchepinsky'
        },
        sidebar: [
            {
                text: 'Documentation',
                items: [
                    { text: 'What is Jellyfin-Newsletter ?', link: '/docs' },
                    { text: 'Installation', link: '/docs/installation' },
                    { text: 'Configuration', link: '/docs/configuration' },
                    { text: 'Placeholders', link: '/docs/placeholders' },
                    { text: 'Themes', link: '/docs/themes' },
                    { text: 'Troubleshooting', link: '/docs/troubleshooting' },
                ]
            }
        ],
        socialLinks: [
            { icon: 'github', link: 'https://github.com/SeaweedbrainCY/jellyfin-newsletter' }
        ],
        outline: {
            level: [2, 3]
        }
    }
})
