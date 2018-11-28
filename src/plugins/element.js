import Vue from 'vue';
import {
  Button,
  Input,
  RadioGroup,
  RadioButton,
  Menu,
  MenuItem
} from 'element-ui';
import lang from 'element-ui/lib/locale/lang/en';
import locale from 'element-ui/lib/locale';

locale.use(lang);

Vue.use(Button);
Vue.use(Input);
Vue.use(RadioGroup);
Vue.use(RadioButton);
Vue.use(Menu);
Vue.use(MenuItem);
