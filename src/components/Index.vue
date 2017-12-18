<template>
    <div class="index-page">
        <form class="main-form" v-if="!loading && !data" @submit.prevent="submit">
            <input
                class="main-form__input"
                type="text"
                name="repo"
                id="repo"
                placeholder="ex: github.com/src-d/go-git"
                v-model="repo"
            >
            <input class="main-form__submit" type="submit" value="Code Burndown!">
        </form>

        <div class="error" v-if="error">{{ error }}</div>

        <div class="loading" v-if="loading">
            <div class="loading__spinner">
                <div class="spinner">loading...</div>
            </div>
            <div class="loading__text">
                Fetching &amp; calculating....
                <br> Please wait, it can take few seconds.
            </div>
        </div>

        <svg width="960" height="500" v-if="data">
          <g transform="translate(50,20)">
            <g class="layer" v-for="(layer, i) in series" :key="i">
              <path class="area" :fill="z(i)" :d="area(layer)"></path>
            </g>
            <g class="axis axis__x" transform="translate(0,450)"></g>
            <g class="axis axis__y"></g>
          </g>
        </svg>
    </div>
</template>

<script>
import * as d3 from 'd3';

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';

export default {
  data() {
    return {
      repo: '',
      loading: false,
      data: null,
      error: null,

      width: 890,
      height: 450
    };
  },

  computed: {
    keys() {
      if (!this.data) {
        return null;
      }

      return d3.range(this.data.data.length);
    },

    series() {
      if (!this.data) {
        return null;
      }

      var stack = d3
        .stack()
        .keys(this.keys)
        .value((d, key) => d[key] || 0);
      return stack(
        this.data.data.map(row => {
          // mutation!!!
          row.total = row.reduce((a, b) => a + b, 0);
          return row;
        })
      );
    },

    x() {
      return d3
        .scaleTime()
        .range([0, this.width])
        .domain([
          new Date(this.data.begin * 1000),
          new Date(this.data.end * 1000)
        ]);
    },

    y() {
      return d3
        .scaleLinear()
        .range([this.height, 0])
        .domain([0, d3.max(this.data.data.map(r => r.total))]);
    },

    z() {
      return d3.scaleOrdinal(d3.schemeCategory20).domain(this.keys);
    },

    step() {
      return (this.data.end - this.data.begin) / (this.data.data.length - 1);
    },

    area() {
      return d3
        .area()
        .x((d, i) => {
          return this.x(new Date((this.data.begin + i * this.step) * 1000));
        })
        .y0(d => this.y(d[0]))
        .y1(d => this.y(d[1]));
    }
  },

  updated() {
    if (!this.data) {
      return;
    }
    this.$nextTick(() => {
      d3.select('.axis__x').call(d3.axisBottom(this.x));
      d3.select('.axis__y').call(d3.axisLeft(this.y));
    });
  },

  methods: {
    submit() {
      this.loading = true;
      this.error = null;

      fetch(`${apiHost}/api/burndown/${this.repo}`)
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
.main-form {
  width: 500px;
  margin: 50px auto;
}

.main-form__input {
  display: block;
  width: 100%;
  margin: 0 auto 10px;
  padding: 10px;

  border: 1px solid #222;
  border-radius: 5px;

  font-size: 24px;
  text-align: center;
}

.main-form__submit {
  display: block;
  margin: 0 auto;
  padding: 10px 20px;

  border: none;
  background: #222;
  color: #fff;
  font-size: 22px;
  text-transform: uppercase;

  -webkit-font-smoothing: antialiased;

  transition: background 150ms;
  cursor: pointer;
}

.main-form__submit:hover {
  background: #90feb5;
}

.loading {
  font-size: 2em;
  text-align: center;
}

.loading__spinner {
  height: 5em;
  padding: 1.5em 0 2.5em;
  margin-bottom: 10px;
}

.spinner,
.spinner:before,
.spinner:after {
  background: #000000;
  animation: spinner-animation 1s infinite ease-in-out;
  width: 1em;
  height: 4em;
}
.spinner {
  padding-top: 2em;
  display: inline-block;
  color: #000000;
  text-indent: -9999em;
  position: relative;
  font-size: 11px;
  transform: translateZ(0);
  animation-delay: -0.16s;
}
.spinner:before,
.spinner:after {
  position: absolute;
  top: 0;
  content: '';
}
.spinner:before {
  left: -1.5em;
  animation-delay: -0.32s;
}
.spinner:after {
  left: 1.5em;
}
@-webkit-keyframes spinner-animation {
  0%,
  80%,
  100% {
    box-shadow: 0 0;
    height: 4em;
  }
  40% {
    box-shadow: 0 -2em;
    height: 5em;
  }
}
@keyframes spinner-animation {
  0%,
  80%,
  100% {
    box-shadow: 0 0;
    height: 4em;
  }
  40% {
    box-shadow: 0 -2em;
    height: 5em;
  }
}
</style>
