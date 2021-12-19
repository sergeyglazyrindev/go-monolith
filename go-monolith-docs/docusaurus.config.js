const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

// With JSDoc @type annotations, IDEs can provide config autocompletion
/** @type {import('@docusaurus/types').DocusaurusConfig} */
(module.exports = {
  title: 'Go-Monolith',
  tagline: 'Build projects in Go easily. Clean code. SOLID. Patterns.',
  url: 'https://gomonolithdocs.sergeyg.me',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'facebook', // Usually your GitHub org/user name.
  projectName: 'docusaurus', // Usually your repo name.

  presets: [
    [
      '@docusaurus/preset-classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/facebook/docusaurus/edit/main/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/facebook/docusaurus/edit/main/website/blog/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'Go-Monolith',
        logo: {
          alt: 'Go-Monolith logo',
          src: 'img/logo.png',
        },
        items: [
          {
            type: 'doc',
            docId: 'intro',
            position: 'left',
            label: 'Tutorial',
          },
          {
            type: 'doc',
            docId: 'demo',
            position: 'left',
            label: 'Demo',
          },
          {
            type: 'doc',
            docId: 'api',
            position: 'left',
            label: 'Api',
          },
          {
            type: 'doc',
            docId: 'contribution',
            position: 'left',
            label: 'Contributing',
          },
          {
            href: 'https://github.com/sergeyglazyrindev/go-monolith',
            label: 'GitHub',
            position: 'right',
          },
          {
            href: 'https://gophers.slack.com/archives/C017ULYJHMZ',
            label: 'Slack',
            position: 'right',
          },
          {
            href: 'https://t.me/joinchat/VzgmokqjF7s4Nzk0',
            label: 'Telegram - uadmin_development',
            position: 'right',
          },
          {
            href: 'https://discord.gg/kADzHWatSj',
            label: 'Discord',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Tutorial',
                to: '/docs/intro',
              },
              {
                label: 'API',
                to: '/docs/api',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Slack',
                href: 'https://gophers.slack.com/archives/C017ULYJHMZ',
              },
              {
                label: 'Telegram - uadmin_development',
                href: 'https://t.me/joinchat/VzgmokqjF7s4Nzk0',
              },
              {
                label: 'Discord',
                href: 'https://discord.gg/kADzHWatSj',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/sergeyglazyrindev/go-monolith',
              },
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/go-monolith',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Go-Monolith, Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      colorMode: {
        disableSwitch: true,
      }
    }),
});
