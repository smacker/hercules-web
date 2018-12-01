<template>
  <div class="page">
    <h-header page="people" :repo="repo" :loading="loading">
      <el-radio-group size="small" v-model="resample" v-if="person">
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

      <div class="content-wrapper" v-if="!loading">
        <div class="sidebar">
          <people :items="people" :selected="person" :onClick="selectPerson"/>
        </div>

        <Responsive v-if="data" class="graph-wrapper">
          <StackGraph
            class="graph"
            slot-scope="props"
            :width="props.width"
            :height="props.height"
            :begin="begin"
            :end="end"
            :data="data"
            :keys="keys"
            :tooltip="!person || resample !== 'raw'"
            :legend="false"
          />
        </Responsive>
      </div>
    </div>
  </div>
</template>


<script>
import Header from "@/components/Header";
import Error from "@/components/Error";
import Loader from "@/components/Loader";
import People from "@/components/People";
import Responsive from "@/components/Responsive";
import StackGraph from "@/components/StackGraph";

import math from "@/lib/math";
import { toMonths, toYears, sumByColumn } from "@/lib/matrix";
import { chooseDefaultResampling } from "@/lib/time";
import differenceInMonths from "date-fns/difference_in_months";
import differenceInYears from "date-fns/difference_in_years";
import { interpolateRdYlBu } from "d3-scale-chromatic";

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || "http://127.0.0.1:8080";

export default {
  props: ["repo"],

  components: {
    HHeader: Header,
    Error,
    Loader,
    People,
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
      person: null,
      resample: "raw"
    };
  },

  computed: {
    resampled() {
      if (!this.person) {
        return null;
      }

      const data = this.serverData[this.person.idx];
      switch (this.resample) {
        case "year":
          return toYears({
            data,
            begin: this.begin,
            end: this.end
          });
        case "month":
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
        const parts = v.split("|");
        const email = parts[parts.length - 1];
        const color = interpolateRdYlBu(i / this.peopleList.length);
        return { value: email, label: v, idx: i, color };
      });
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
      this.loading = true;
      this.error = null;

      fetch(`${apiHost}/api/burndown/${this.repo}`)
        .then(r => r.json())
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          if (json.data.peopleData.length < 2) {
            return Promise.reject("Not enough data");
          }

          this.serverData = json.data.peopleData;
          this.peopleList = json.data.peopleList;
          this.overallData = math.transpose(
            json.data.peopleList.map((_, i) => {
              return sumByColumn(json.data.peopleData["" + i]);
            })
          );
          this.begin = json.data.begin;
          this.end = json.data.end;
          this.resample = chooseDefaultResampling(this.begin, this.end);
        })
        .catch(e => (this.error = e))
        .then(() => (this.loading = false));
    },

    selectPerson(person) {
      this.person = person;
    }
  }
};
</script>

<style scoped>
.content-wrapper {
  display: flex;
  height: 100%;
}

.sidebar {
  flex: 0 0 auto;
  width: 400px;
  height: 100%;
  padding-bottom: 10px;
}

.graph-wrapper {
  flex: 1 1 100%;

  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.graph {
  margin: 0 auto;
}
</style>
