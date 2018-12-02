<template>
  <div class="page">
    <h-header page="overall" :repo="repo" :loading="loading">
      <el-radio-group size="small" v-model="resample">
        <el-radio-button
          v-for="opt in resampleOptions"
          :key="opt.name"
          :label="opt.name"
          :disabled="opt.disabled"
        />
      </el-radio-group>
    </h-header>

    <error :msg="error" v-if="error"/>

    <div class="page-body" v-if="!error">
      <loader v-if="loading"/>

      <Responsive v-if="data" class="graph-wrapper">
        <StackGraph
          class="graph"
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
</template>


<script>
import Header from "@/components/Header";
import Error from "@/components/Error";
import Loader from "@/components/Loader";
import Responsive from "@/components/Responsive";
import StackGraph from "@/components/StackGraph";

import { fetch } from "@/lib/api";
import math from "@/lib/math";
import { toMonths, toYears } from "@/lib/matrix";
import { chooseDefaultResampling } from "@/lib/time";
import differenceInMonths from "date-fns/difference_in_months";
import differenceInYears from "date-fns/difference_in_years";

export default {
  props: ["repo"],

  components: {
    HHeader: Header,
    Error,
    Loader,
    Responsive,
    StackGraph
  },

  data() {
    return {
      loading: true,
      resample: "raw",
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
        case "raw":
          return {
            data: this.serverData,
            keys: math.range(0, this.serverData.length).toArray()
          };

        case "month":
          return toMonths({
            data: this.serverData,
            begin: this.begin,
            end: this.end
          });

        case "year":
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
        { name: "raw", disabled: false },
        { name: "month", disabled: !totalMonths || totalMonths > 50 },
        { name: "year", disabled: !totalYears || totalYears == 1 }
      ];
    }
  },

  created() {
    this.fetchData();
  },

  watch: {
    $route: "fetchData"
  },

  methods: {
    fetchData() {
      this.serverData = null;
      this.loading = true;
      this.error = null;

      fetch(`/api/analysis/project/${this.repo}`)
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          if (json.project.length < 2) {
            return Promise.reject("Not enough data");
          }

          this.serverData = json.project;
          this.begin = json.begin;
          this.end = json.end;
          this.resample = chooseDefaultResampling(this.begin, this.end);
        })
        .catch(e => (this.error = e))
        .then(() => (this.loading = false));
    }
  }
};
</script>

<style scoped>
.graph-wrapper {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.graph {
  margin: 0 auto;
}
</style>
