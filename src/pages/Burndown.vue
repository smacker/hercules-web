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

        <span>
          Resample: <select v-model="resample">
            <option v-for="opt in resampleOptions" :key="opt.name" :disabled="opt.disabled">{{opt.name}}</option>
          </select>
        </span>

        <Responsive v-if="data" class="graph">
          <StackGraph
            slot-scope="props"
            :width="props.width"
            :height="props.height"
            :begin="begin"
            :end="end"
            :data="data.data"
            :keys="data.keys"
            :tooltip="resample != 'raw'"
            :legend="resample != 'raw' && data.keys.length < 10"
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
import { toMonths, toYears } from '@/lib/matrix';
import differenceInMonths from 'date-fns/difference_in_months';
import differenceInYears from 'date-fns/difference_in_years';

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
      resample: 'raw',
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

      switch (this.resample) {
        case 'raw':
          return {
            data: this.serverData,
            keys: math.range(0, this.serverData.length).toArray()
          };

        case 'month':
          return toMonths({
            data: this.serverData,
            begin: this.begin,
            end: this.end
          });

        case 'year':
          return toYears({
            data: this.serverData,
            begin: this.begin,
            end: this.end
          });
        default:
          return null;
      }
    },

    resampleOptions() {
      let totalMonths = 0;
      let totalYears = 0;
      if (this.end && this.begin) {
        const begin = new Date(this.begin * 1000);
        const end = new Date(this.end * 1000);
        totalMonths = differenceInMonths(end, begin);
        totalYears = differenceInYears(end, begin);
      }

      return [
        { name: 'raw', disabled: false },
        { name: 'month', disabled: !totalMonths || totalMonths > 50 },
        { name: 'year', disabled: !totalYears || totalYears == 1 }
      ];
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
          if (json.data.data.length < 2) {
            return Promise.reject('Not enough data');
          }

          this.serverData = json.data.data;
          this.begin = json.data.begin;
          this.end = json.data.end;
          this.resample = chooseDefaultResampling(this.begin, this.end);
        })
        .catch(e => (this.error = e))
        .then(() => (this.loading = false));
    }
  }
};

function chooseDefaultResampling(begin, end) {
  begin = new Date(begin * 1000);
  end = new Date(end * 1000);

  const years = differenceInYears(end, begin);
  if (years >= 3) {
    return 'year';
  }
  const months = differenceInMonths(end, begin);
  if (months > 3 && months < 36) {
    return 'month';
  }
  return 'raw';
}
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
