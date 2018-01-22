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
        <router-link :to="`/${repo}/burndown/people`">By people</router-link>

        <div class="wrapper">
          <files-tree
            class="sidebar"
            :tree="filesTree.children"
            :list="filesList"
            :onSelect="selectFile"
          />

          <div class="content">
            <div class="current-file">{{currentFile.path}}</div>
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
    </div>
  </div>
</template>

<script>
import Spinner from '@/components/Spinner';
import Responsive from '@/components/Responsive';
import StackGraph from '@/components/StackGraph';

import math from 'mathjs';
import { toMonths, toYears } from '@/lib/matrix';
import { chooseDefaultResampling } from '@/lib/time';

import { filesToTree } from '@/lib/files';
import FilesTree from '@/components/FilesTree';

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';
const initialState = {
  loading: true,
  error: null,
  filesTree: null,
  filesList: null,
  currentFile: null,

  serverData: null,
  begin: null,
  end: null,
  resample: 'raw'
};

function resetState(instance) {
  Object.keys(initialState).forEach(key => {
    instance[key] = initialState[key];
  });
}

export default {
  props: ['repo'],

  components: {
    Spinner,
    Responsive,
    StackGraph,
    FilesTree
  },

  data() {
    return initialState;
  },

  computed: {
    data() {
      if (!this.serverData) {
        return null;
      }

      if (!this.currentFile) {
        return null;
      }

      const data = this.serverData[this.currentFile.path];

      switch (this.resample) {
        case 'raw':
          return {
            data,
            keys: math.range(0, data.length).toArray()
          };

        case 'month':
          return toMonths({
            data,
            begin: this.begin,
            end: this.end
          });

        case 'year':
          return toYears({
            data,
            begin: this.begin,
            end: this.end
          });
        default:
          return null;
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
      resetState(this);

      fetch(`${apiHost}/api/burndown/${this.repo}`)
        .then(r => r.json())
        .then(json => {
          if (json.error) {
            return Promise.reject(json.error);
          }
          const { tree, list } = filesToTree(Object.keys(json.data.filesData));
          this.filesTree = tree;
          this.filesList = list;
          this.serverData = json.data.filesData;
          this.begin = json.data.begin;
          this.end = json.data.end;
          this.resample = chooseDefaultResampling(this.begin, this.end);

          this.currentFile = list[0];
        })
        .catch(e => (this.error = e))
        .then(() => (this.loading = false));
    },

    selectFile(file) {
      this.currentFile = file;
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

.wrapper {
  display: flex;
}

.sidebar {
  flex: 0 0 auto;
  width: 400px;
}

.content {
  width: 100%;
}

.current-file {
  text-align: center;
}
</style>
