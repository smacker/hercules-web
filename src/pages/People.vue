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
      <loader v-if="loading" :status="loadingStatus"/>

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
import analysisMixin from "./analysisMixin";

import { toMonths, toYears, sumByColumn, transpose } from "@/lib/matrix";
import { chooseDefaultResampling } from "@/lib/time";
import differenceInMonths from "date-fns/difference_in_months";
import differenceInYears from "date-fns/difference_in_years";
import { interpolateRdYlBu } from "d3-scale-chromatic";

export default {
  mixins: [analysisMixin],
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
      serverData: null,
      allPeopleList: [],
      peopleList: [],
      overallData: null,
      begin: null,
      end: null,
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
            keys: Array.from(Array(data.length).keys())
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
      return this.peopleList
        .map((v, i) => {
          const idx = this.allPeopleList.indexOf(v);
          const color = interpolateRdYlBu(i / this.peopleList.length);
          return { value: v, label: v, idx, color };
        })
        .filter(v => !!v);
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
      this.fetch(`/api/analysis/people/${this.repo}`, json => {
        if (json.peopleData.length < 2) {
          return Promise.reject("Not enough data");
        }

        this.serverData = json.peopleData;
        this.allPeopleList = json.peopleList;
        this.peopleList = json.peopleList;
        this.overallData = transpose(
          json.peopleList.map((_, i) => {
            return sumByColumn(json.peopleData["" + i]);
          })
        );
        this.begin = json.begin;
        this.end = json.end;
        this.resample = chooseDefaultResampling(this.begin, this.end);

        this.groupPeople();
        this.sortPeople();
      });
    },

    groupPeople() {
      // minimum number of people to apply grouping
      const minPeople = 10;
      if (this.peopleList.length <= minPeople) {
        return;
      }

      // max total overall value
      const max = this.overallData.reduce((acc, col) => {
        const sum = col.reduce((a, b) => a + b, 0);
        if (sum > acc) {
          return sum;
        }
        return acc;
      }, 0);

      const threshold = max * 0.01;
      // only people that contributed more than threshold
      const keep = this.overallData.reduce((acc, col) => {
        return acc.concat(
          col.reduce((innerAcc, v, i) => {
            if (!acc.includes(i) && v > threshold) {
              return innerAcc.concat([i]);
            }
            return innerAcc;
          }, [])
        );
      }, []);

      // drop grouping if there are too few survived people
      if (keep.length < minPeople) {
        return;
      }

      const overallData = this.overallData.map(col => {
        let restVal = 0;

        return col
          .filter((v, i) => {
            if (keep.includes(i)) {
              return true;
            }

            restVal += v;
            return false;
          })
          .concat([restVal]);
      });
      const people = this.peopleList
        .filter((_, i) => keep.includes(i))
        .concat(["Others"]);

      this.peopleList = people;
      this.overallData = overallData;
    },

    sortPeople() {
      // sort by the value at the end of timeline
      const endOverallData = this.overallData[this.overallData.length - 1];
      const mapped = endOverallData.map((v, i) => ({ v, i }));
      mapped.sort((a, b) => b.v - a.v);
      // resort people and data accordingly
      this.peopleList = mapped.map(v => this.peopleList[v.i]);
      this.overallData = this.overallData.map(col => {
        return mapped.map(v => col[v.i]);
      });
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
