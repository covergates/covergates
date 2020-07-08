import _Vue from 'vue';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';

function highlight(text: string, lang?: string): string {
  if (lang) {
    return hljs.highlight(lang, text).value;
  }
  return hljs.highlightAuto(text).value;
}

export function HighlightPlugin(Vue: typeof _Vue): void {
  hljs.configure({
    tabReplace: '    '
  });
  Vue.prototype.$highlight = highlight;
}

declare module 'vue/types/vue' {
  interface Vue {
    readonly $highlight: typeof highlight;
  }
}
