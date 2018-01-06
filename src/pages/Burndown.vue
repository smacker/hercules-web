<template>
  <div class="burndown-page">
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

    <div>
      Mode: <select v-model="mode">
        <option>raw</option>
        <option>year</option>
      </select>
      </div>

    <Responsive v-if="data" class="graph">
      <Burndown
        slot-scope="props"
        :width="props.width"
        :height="props.height"
        :begin="begin"
        :end="end"
        :data="data.data"
        :keys="data.keys"
      />
    </Responsive>
  </div>
</template>


<script>
import Spinner from '@/components/Spinner';
import Responsive from '@/components/Responsive';
import Burndown from '@/components/Burndown';

import math from 'mathjs';
import { toYear } from '@/lib/matrix';

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';

export default {
  props: ['repo'],

  components: {
    Spinner,
    Responsive,
    Burndown
  },

  data() {
    return {
      loading: true,
      mode: 'raw',
      serverData: null,
      begin: null,
      end: null,
      error: null
    };
  },

  computed: {
    data() {
      if (!this.serverData) {
        return null;
      }

      switch (this.mode) {
        case 'raw':
          return {
            data: this.serverData,
            keys: math.range(0, this.serverData.length).toArray()
          };

        default:
          const { keys, matrix } = toYear({
            data: this.serverData,
            begin: this.begin,
            end: this.end
          });

          return {
            data: math.transpose(matrix).toArray(),
            keys
          };
      }
    }
  },

  created() {
    this.fetchData();
  },

  watch: {
    $route: 'fetchData'
  },

  methods: {
    fetchData() {
      this.serverData = null;
      this.loading = true;
      this.error = null;

      fetch(`${apiHost}/api/burndown/${this.repo}`)
        .then(r => r.json())
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          this.serverData = json.data.data;
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
