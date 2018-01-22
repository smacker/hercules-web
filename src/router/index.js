import Vue from 'vue';
import Router from 'vue-router';
import Index from '@/pages/Index';
import Project from '@/pages/Project';
import People from '@/pages/People';
import Files from '@/pages/Files';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Heracules',
      component: Index
    },
    {
      path: '/:repo(.+)/burndown/people',
      component: People,
      props: true
    },
    {
      path: '/:repo(.+)/burndown/files',
      component: Files,
      props: true
    },
    {
      path: '/:repo(.+)/burndown',
      component: Project,
      props: true
    }
  ]
});
