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
        <v-select v-model="person" :options="people" class="v-select"></v-select>

        <Responsive v-if="data" class="graph">
          <StackGraph
            slot-scope="props"
            :width="props.width"
            :height="props.height"
            :begin="begin"
            :end="end"
            :data="data"
            :keys="keys"
            :tooltip="!person || resampling !== 'raw'"
            :legend="!person || resampling === 'year'"
          />
        </Responsive>
      </div>
    </div>
  </div>
</template>


<script>
import vSelect from 'vue-select';
import Spinner from '@/components/Spinner';
import Responsive from '@/components/Responsive';
import StackGraph from '@/components/StackGraph';

import math from 'mathjs';
import { toMonths, toYears, sumByColumn } from '@/lib/matrix';
import { chooseDefaultResampling } from '@/lib/time';

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';

export default {
  props: ['repo'],

  components: {
    vSelect,
    Spinner,
    Responsive,
    StackGraph
  },

  data() {
    return {
      loading: true,
      serverData: null,
      peopleList: [],
      overallData: null,
      begin: null,
      end: null,
      error: null,
      person: null
    };
  },

  computed: {
    resampling() {
      if (!this.person) {
        return null;
      }
      return chooseDefaultResampling(this.begin, this.end);
    },

    resampled() {
      if (!this.person) {
        return null;
      }

      const data = this.serverData[this.person.idx];
      switch (this.resampling) {
        case 'year':
          return toYears({
            data,
            begin: this.begin,
            end: this.end
          });
        case 'month':
          return toMonths({
            data,
            begin: this.begin,
            end: this.end
          });
        default:
          return {
            data,
            keys: math.range(0, data.length).toArray()
          };
      }
    },

    data() {
      if (this.person) {
        return this.resampled.data;
      }
      return this.overallData;
    },

    keys() {
      if (this.person) {
        return this.resampled.keys;
      }
      return this.peopleList;
    },

    people() {
      return this.peopleList.map((v, i) => {
        const parts = v.split('|');
        const email = parts[parts.length - 1];
        return { value: email, label: v, idx: i };
      });
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
      this.loading = true;
      this.error = null;

      fetch(`${apiHost}/api/burndown/${this.repo}`)
        .then(r => r.json())
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          if (json.data.peopleData.length < 2) {
            return Promise.reject('Not enough data');
          }

          this.serverData = json.data.peopleData;
          this.peopleList = json.data.peopleList;
          this.overallData = math.transpose(
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

.v-select {
  display: inline-block;
  width: 700px;
}
</style>
