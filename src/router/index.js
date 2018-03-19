import Vue from 'vue';
import Router from 'vue-router';
import Index from '@/pages/Index';
import Project from '@/pages/Project';
import People from '@/pages/People';
import Files from '@/pages/Files';

Vue.use(Router);

const router = new Router({
  routes: [
    {
      path: '/',
      name: 'Heracules',
      component: Index
    },
    {
      path: '/:repo(.+)/burndown/people',
      component: People,
      name: 'people',
      props: true
    },
    {
      path: '/:repo(.+)/burndown/files',
      component: Files,
      name: 'files',
      props: true
    },
    {
      path: '/:repo(.+)/burndown',
      component: Project,
      name: 'project',
      props: true
    }
  ]
});

// Monkey patch all the things!
// actually I just want to use named routes
// and keep '/' in the url instead of ugly %2F
const originalMatch = router.match.bind(router);
router.match = (raw, current, redirectFrom) => {
  let val = originalMatch(raw, current, redirectFrom);
  if (val && val.fullPath) {
    val = Object.assign({}, val, {
      fullPath: val.fullPath.replace(/%2F/g, '/')
    });
    Object.freeze(val);
  }
  return val;
};

export default router;
