<template>
  <div>
    <div class="search">
      <input type="text" placeholder="start typing file name" v-model="query" />
    </div>
    <div class="listing">
      <file-item
        v-if="!query"
        v-for="item in tree"
        :key="item.path"
        :item="item"
        :onSelect="onSelect"
      />
      <div
        v-if="query"
        v-for="item in filtered"
        :key="item.path"
        :class="{'file-item': true, [item.type]:true}"
        @click="onSelect(item)"
      >
        {{item.path}}
      </div>
    </div>
  </div>
</template>

<script>
import FileItem from '@/components/FileItem';

export default {
  props: ['tree', 'list', 'onSelect'],

  components: {
    FileItem
  },

  data() {
    return {
      query: ''
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
.search input {
  display: block;
  width: 100%;
  padding: 5px 10px;
  margin-bottom: 3px;
  border: 1px solid #eaecef;
}

.listing {
  max-height: 511px;
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
