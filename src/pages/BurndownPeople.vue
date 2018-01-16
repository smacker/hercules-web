<template>
  <div class="burndown-page">
    <div class="error" v-if="error">
      <h3>Oops! There is an error:</h3>
      <p>{{ error }}</p>
      <p><a href="/">Try another repository</a></p>
    </div>

    <div v-if="!error">
      <div class="loading" v-if="loading">
        <div class="loading__spinner">
          <spinner />
        </div>
        <div class="loading__text">
          Fetching &amp; calculating....
          <br> Please wait, it can take few seconds.
        </div>
      </div>

      <div v-if="!loading">
        <a href="/">Back</a>
        <router-link :to="`/${repo}/burndown`">Project overall</router-link>

        <Responsive v-if="data" class="graph">
          <StackGraph
            slot-scope="props"
            :width="props.width"
            :height="props.height"
            :begin="begin"
            :end="end"
            :data="data"
            :keys="keys"
            :tooltip="true"
            :legend="true"
          />
        </Responsive>
      </div>
    </div>
  </div>
</template>


<script>
import Spinner from '@/components/Spinner';
import Responsive from '@/components/Responsive';
import StackGraph from '@/components/StackGraph';

import math from 'mathjs';
import { sumByColumn } from '@/lib/matrix';

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';

export default {
  props: ['repo'],

  components: {
    Spinner,
    Responsive,
    StackGraph
  },

  data() {
    return {
      loading: true,
      data: null,
      begin: null,
      end: null,
      keys: null,
      error: null
    };
  },

  created() {
    this.fetchData();
  },

  watch: {
    $route: 'fetchData'
  },

  methods: {
    fetchData() {
      this.loading = true;
      this.error = null;

      fetch(`${apiHost}/api/burndown/${this.repo}`)
        .then(r => r.json())
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          if (json.data.project.length < 2) {
            return Promise.reject('Not enough data');
          }

          this.keys = json.data.peopleList;
          this.data = math.transpose(
            json.data.peopleList.map((_, i) => {
              return sumByColumn(json.data.peopleData['' + i]);
            })
          );
          this.begin = json.data.begin;
          this.end = json.data.end;
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

.graph {
  min-width: 600px;
  max-width: 1200px;
  margin: 0 auto;
}

.error {
  color: #f5222d;
}
</style>
