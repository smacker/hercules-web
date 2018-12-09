<template>
  <div class="graph" :style="{width: `${this.width}px`, height: `${this.height}px`,}">
    <StackSVG
      :begin="begin"
      :end="end"
      :data="data"
      :width="width"
      :height="height"
      :keys="keys"
      :colors="colors"
      @layerMouseover="showTooltip"
      @layerMousemove="moveTooltip"
      @layerMouseout="hideTooltip"
    />
    <Legend :colors="colors" :items="keys" v-if="legend"/>
    <div class="tooltip" v-if="tooltip && tooltipShown" :style="tooltipPosition">{{tooltipContent}}</div>
  </div>
</template>

<script>
import { interpolateRdYlBu } from "d3-scale-chromatic";
import StackSVG from "@/components/StackSVG";
import Legend from "@/components/Legend";

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
    Legend,
    StackSVG
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
    this.updatePosition();
  },

  updated() {
    this.updatePosition();
  },

  computed: {
    colors() {
      return Array.from(Array(this.keys.length).keys()).map(v =>
        this.colorSchema(v / this.keys.length)
      );
    }
  },

  methods: {
    updatePosition() {
      const rect = this.$el.getBoundingClientRect();
      this.top = rect.top;
      this.left = rect.left;
    },

    showTooltip(layerId, e) {
      this.tooltipContent = this.keys[layerId];
      this.moveTooltip(layerId, e);
      this.tooltipShown = true;
    },
    moveTooltip(layerId, e) {
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
