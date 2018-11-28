<template>
  <div class="page">
    <h-header page="files" :repo="repo" :loading="loading"/>

    <error :msg="error" v-if="error"/>

    <div class="page-body" v-if="!error">
      <loader v-if="loading"/>

      <div class="content-wrapper" v-if="!loading">
        <files-tree
          class="sidebar"
          :tree="filesTree.children"
          :list="filesList"
          :onSelect="selectFile"
        />

        <div class="content">
          <div class="current-file">{{currentFile.path}}</div>
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
    </div>
  </div>
</template>

<script>
import Header from "@/components/Header";
import Error from "@/components/Error";
import Loader from "@/components/Loader";
import Responsive from "@/components/Responsive";
import StackGraph from "@/components/StackGraph";

import math from "@/lib/math";
import { toMonths, toYears } from "@/lib/matrix";
import { chooseDefaultResampling } from "@/lib/time";

import { filesToTree } from "@/lib/files";
import FilesTree from "@/components/FilesTree";

const hercules = window.hercules || {};
const apiHost = hercules.apiHost || "http://127.0.0.1:8080";
const initialState = {
  loading: true,
  error: null,
  filesTree: null,
  filesList: null,
  currentFile: null,

  serverData: null,
  begin: null,
  end: null,
  resample: "raw"
};

function resetState(instance) {
  Object.keys(initialState).forEach(key => {
    instance[key] = initialState[key];
  });
}

export default {
  props: ["repo"],

  components: {
    HHeader: Header,
    Error,
    Loader,
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
        case "raw":
          return {
            data,
            keys: math.range(0, data.length).toArray()
          };

        case "month":
          return toMonths({
            data,
            begin: this.begin,
            end: this.end
          });

        case "year":
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
    $route: "fetchData"
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

.content {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.current-file {
  flex: 0 0 auto;

  height: 28px;
  line-height: 28px;
  margin-bottom: 3px;
  text-align: center;
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
