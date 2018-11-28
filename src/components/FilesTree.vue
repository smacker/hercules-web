<template>
  <div class="files-tree">
    <div class="search">
      <el-input
        placeholder="start typing file name"
        v-model="query"
        size="mini"
        :clearable="true"
        class="input"
      />
    </div>
    <div class="listing">
      <div v-if="!query">
        <file-item v-for="item in tree" :key="item.path" :item="item" :onSelect="onSelect"/>
      </div>
      <div v-if="query">
        <div
          v-for="item in filtered"
          :key="item.path"
          :class="{'file-item': true, [item.type]:true}"
          @click="onSelect(item)"
        >{{item.path}}</div>
      </div>
    </div>
  </div>
</template>

<script>
import FileItem from "@/components/FileItem";

export default {
  props: ["tree", "list", "onSelect"],

  components: {
    FileItem
  },

  data() {
    return {
      query: ""
    };
  },

  computed: {
    filtered() {
      if (!this.query) {
        return [];
      }

      return this.list.filter(i => i.path.includes(this.query));
    }
  },

  methods: {}
};
</script>

<style scoped>
.files-tree {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.search {
  flex: 0 0 auto;
}

.input {
  margin-bottom: 3px;
}

.listing {
  flex: 0 1 auto;
  overflow: auto;

  background: #fff;
  border: 1px solid #eaecef;
}

.file-item {
  padding: 5px 10px;
  border-bottom: 1px solid #eaecef;
  cursor: pointer;
}

.file-item:hover {
  background: #f6f8fa;
}
</style>
