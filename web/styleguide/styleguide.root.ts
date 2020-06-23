import vuetify from '../src/plugins/vuetify';
import Vue from 'vue';
import Component from 'vue-class-component';

@Component({
  template: '<div class="v-application v-application--is-ltr theme--light"><slot/></div>'
})
class App extends Vue {
  // type inference enabled
}

export default (previewComponent: Vue.Component) => {
  return {
    el: '#app',
    vuetify,
    render(createElement: Vue.CreateElement) {
      return createElement(App, [createElement(previewComponent)]);
    }
  };
};
