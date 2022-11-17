export default [
  {
    path: '/account',
    routes: [
      {
        path: '/account',
        routes: [
          {
            name: 'login',
            path: '/account/login',
            layout: false,
            component: './Account/Login',
          },
          {
            name: 'settings',
            path: '/account/settings',
            component: './Account/Setting',
          },
        ],
      },
      {
        component: './404',
      },
    ],
  },
  {
    name: 'namespace',
    icon: 'database',
    path: '/app/namespace',
    component: '@/pages/App/Namespace',
  },
  {
    name: 'action',
    icon: 'code',
    path: '/app/action',
    routes: [
      {
        name: 'backend',
        path: '/app/action/backend',
        component: '@/pages/App/Action/Backend',
      },
      {
        name: 'static',
        path: '/app/action/static',
        component: '@/pages/App/Action/Static',
      },
    ],
  },
  {
    name: 'site',
    icon: 'cluster',
    path: '/app/site',
    component: '@/pages/App/Site',
  },
  {
    path: '/',
    redirect: '/app/namespace',
  },
  {
    component: './404',
  },
];
