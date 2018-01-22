export function filesToTree(files) {
  const list = [];
  const tree = files.reduce(
    (acc, f) => {
      const path = f.split('/');
      const dirs = path.slice(0, path.length - 1);
      const file = path[path.length - 1];
      let cur = acc;
      dirs.forEach(d => {
        let item = cur.children.find(i => i.name === d);
        if (!item) {
          item = {
            name: d,
            type: 'dir',
            level: (cur.level || 0) + 1,
            path: cur.path ? [cur.path, d].join('/') : d,
            children: [],
            open: false
          };
          cur.children.push(item);
        }
        cur = item;
      });
      const fileObj = {
        name: file,
        type: 'file',
        level: dirs.length + 1,
        path: f
      };

      list.push(fileObj);
      cur.children.push(fileObj);

      return acc;
    },
    { name: 'root', children: [] }
  );

  childrenSort(tree);
  return { tree, list };
}

function childrenSort(item) {
  if (!item.children) {
    return;
  }

  item.children.sort((a, b) => {
    // dirs on top
    if (a.type === 'file' && b.type === 'dir') {
      return 1;
    }
    if (a.type === 'dir' && b.type === 'file') {
      return -1;
    }
    // by name otherwise
    if (a.name > b.name) {
      return 1;
    }
    if (a.name < b.name) {
      return -1;
    }
    return 0;
  });

  item.children.forEach(child => childrenSort(child));
}
