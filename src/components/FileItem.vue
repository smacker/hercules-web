<template>
  <div>
    <div
      :class="{'file-item': true, [item.type]:true}"
      :style="{paddingLeft: item.level*10+'px'}"
      @click="item.type === 'dir' ? onClick(item) : onSelect(item)"
    >
      <icon name="folder" v-if="item.type === 'dir'"/>
      <icon name="file-code" v-if="item.type === 'file'"/>
      {{item.name}}
    </div>
    <div v-if="item.type === 'dir' && item.open">
      <file-item
        v-for="child in item.children"
        :key="child.path"
        :item="child"
        :onSelect="onSelect"
      />
    </div>
  </div>
</template>

<script>
import "vue-awesome/icons/folder";
import "vue-awesome/icons/file-code";
import Icon from "vue-awesome/components/Icon";

export default {
  props: ["item", "onSelect"],
  name: "file-item",

  components: {
    Icon
  },

  methods: {
    onClick(item) {
      item.open = !item.open;
    }
  }
};
</script>

<style scoped>
/* FIXME: it's copy-past from FilesTree, need it because of local scope */

.file-item {
  padding: 5px 10px;
  border-bottom: 1px solid #eaecef;
  cursor: pointer;
}

.file-item:hover {
  background: #f6f8fa;
}
</style>
