import Vue from 'vue';
import Router from 'vue-router';
import Index from '@/pages/Index';
import Burndown from '@/pages/Burndown';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Heracules',
      component: Index
    },
    {
      path: '/:repo(.+)/burndown',
      component: Burndown
    }
  ]
});
