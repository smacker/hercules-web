<template>
  <div class="graph" :style="{width: `${this.width}px`, height: `${this.height}px`,}">
    <svg :width="this.width" :height="this.height">
      <g :transform="gTransform">
        <g
          class="layer"
          v-for="(layer, i) in series"
          :key="i"
          @mouseover="showTooltip(i, $event)"
          @mousemove="moveTooltip"
          @mouseout="hideTooltip"
        >
          <path class="area" :fill="z(i)" :d="area(layer)"></path>
        </g>
        <g class="axis axis__x" :transform="xAxisTransform"></g>
        <g class="axis axis__y"></g>
      </g>
    </svg>
    <Legend :colors="colors" :items="keys" v-if="legend"/>
    <div class="tooltip" v-if="tooltip && tooltipShown" :style="tooltipPosition">{{tooltipContent}}</div>
  </div>
</template>

<script>
import { stack, area } from "d3-shape";
import { range, max } from "d3-array";
import { select } from "d3-selection";
import { scaleTime, scaleLinear, scaleOrdinal } from "d3-scale";
import { axisLeft, axisBottom } from "d3-axis";
import { interpolateRdYlBu } from "d3-scale-chromatic";

import Legend from "@/components/Legend";

const d3 = {
  range,
  stack,
  area,
  max,
  scaleTime,
  scaleLinear,
  scaleOrdinal,
  select,
  axisBottom,
  axisLeft
};

const margin = { top: 5, right: 20, bottom: 30, left: 50 };

export default {
  props: {
    begin: {
      type: Number,
      required: true
    },
    end: {
      type: Number,
      required: true
    },
    data: {
      type: Array,
      required: true
    },
    width: {
      type: Number,
      default: 960
    },
    height: {
      type: Number,
      default: 500
    },
    keys: {
      type: Array,
      required: true
    },
    colorSchema: {
      type: Function,
      default: interpolateRdYlBu
    },
    legend: {
      type: Boolean,
      default: true
    },
    tooltip: {
      type: Boolean,
      default: true
    }
  },

  components: {
    Legend
  },

  data() {
    return {
      top: 0,
      left: 0,

      tooltipShown: false,
      tooltipPosition: {},
      tooltipContent: ""
    };
  },

  mounted() {
    this.updateAxis();
    this.updatePosition();
  },

  updated() {
    this.updateAxis();
    this.updatePosition();
  },

  computed: {
    gWidth() {
      return this.width - margin.left - margin.right;
    },

    gHeight() {
      return this.height - margin.top - margin.bottom;
    },

    gTransform() {
      return `translate(${margin.left},${margin.top})`;
    },

    colors() {
      return d3
        .range(0, this.keys.length)
        .map(v => this.colorSchema(v / this.keys.length));
    },

    // graph
    series() {
      const keyIdx = this.keys.reduce((r, k, i) => {
        r[k] = i;
        return r;
      }, {});
      var stack = d3
        .stack()
        .keys(this.keys)
        .value((d, key) => d[keyIdx[key]] || 0);
      return stack(
        this.data.map(row => {
          // mutation!!!
          row.total = row.reduce((a, b) => a + b, 0);
          return row;
        })
      );
    },

    x() {
      return d3
        .scaleTime()
        .range([0, this.gWidth])
        .domain([new Date(this.begin * 1000), new Date(this.end * 1000)]);
    },

    y() {
      return d3
        .scaleLinear()
        .range([this.gHeight, 0])
        .domain([0, d3.max(this.data.map(r => r.total))]);
    },

    z() {
      return d3.scaleOrdinal(this.colors).domain(this.keys);
    },

    step() {
      return (this.end - this.begin) / (this.data.length - 1);
    },

    area() {
      return d3
        .area()
        .x((d, i) => {
          return this.x(new Date((this.begin + i * this.step) * 1000));
        })
        .y0(d => this.y(d[0]))
        .y1(d => this.y(d[1]));
    },

    // axis
    xAxisTransform() {
      return `translate(0,${this.gHeight})`;
    }
  },

  methods: {
    updatePosition() {
      const rect = this.$el.getBoundingClientRect();
      this.top = rect.top;
      this.left = rect.left;
    },

    updateAxis() {
      this.$nextTick(() => {
        d3.select(".axis__x").call(d3.axisBottom(this.x));
        d3.select(".axis__y").call(d3.axisLeft(this.y));
      });
    },

    showTooltip(layerId, e) {
      this.tooltipContent = this.keys[layerId];
      this.moveTooltip(e);
      this.tooltipShown = true;
    },
    moveTooltip(e) {
      this.tooltipPosition = {
        top: e.clientY + 5 - this.top + "px",
        left: e.clientX + 5 - this.left + "px"
      };
    },
    hideTooltip() {
      this.tooltipShown = false;
    }
  }
};
</script>


<style scoped>
.graph {
  position: relative;
}

.legend {
  position: absolute;
  top: 20px;
  left: 60px;
}

.tooltip {
  position: absolute;
  border: 1px solid #333;
  background: #fff;
  padding: 10px;
}
</style>
