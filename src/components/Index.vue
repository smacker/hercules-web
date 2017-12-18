<template>
    <div class="index-page">
        <main-form :on-submit="submit" v-if="!loading && !data" />
        <div class="error" v-if="error">{{ error }}</div>

        <div class="loading" v-if="loading">
            <div class="loading__spinner">
                <spinner />
            </div>
            <div class="loading__text">
                Fetching &amp; calculating....
                <br> Please wait, it can take few seconds.
            </div>
        </div>

        <Burndown v-if="data" :begin="data.begin" :end="data.end" :data="data.data" />
    </div>
</template>

<script>
import MainForm from './MainForm';
import Spinner from './Spinner';
import Burndown from './Burndown';

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';

export default {
  components: {
    MainForm,
    Spinner,
    Burndown
  },

  data() {
    return {
      loading: false,
      data: null,
      error: null
    };
  },

  methods: {
    submit(repo) {
      this.loading = true;
      this.error = null;

      fetch(`${apiHost}/api/burndown/${repo}`)
        .then(r => r.json())
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          this.data = json.data;
        })
        .catch(e => (this.error = e))
        .then(() => (this.loading = false));
    }
  }
};
</script>

<style scoped>
.loading {
  font-size: 2em;
  text-align: center;
}

.loading__spinner {
  height: 5em;
  padding: 1.5em 0 2.5em;
  margin-bottom: 10px;
}

.error {
  position: absolute;
  bottom: 0;
  right: 0;

  display: inline-block;
  margin: 0 10px 10px 0;
  padding: 10px;

  background: #fff1f0;
  border: 1px solid #ffa39e;
  color: #f5222d;
}
</style>
