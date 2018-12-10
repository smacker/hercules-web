import { fetch } from '@/lib/api';

export default {
  data() {
    return {
      loading: true,
      loadingStatus: null,
      error: null
    };
  },

  methods: {
    fetch(endpoint, callback) {
      this.loading = true;
      this.loadingStatus = null;
      this.error = null;

      const call = () =>
        fetch(endpoint)
          .then(json => {
            if (json.error) {
              return Promise.reject(json.error);
            }

            return json;
          })
          .then(json => {
            if (typeof json.status !== 'undefined') {
              this.loadingStatus = json.status;

              return new Promise(resolve => {
                setTimeout(() => resolve(call()), 2000);
              });
            }

            return callback(json);
          });

      return call()
        .catch(e => (this.error = e))
        .then(() => (this.loading = false));
    }
  }
};
