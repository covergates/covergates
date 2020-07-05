declare module 'vue-countup-v2' {
  export interface ICountUp { /* eslint-disable-line */
    delay: number;
    endVal: number;
    options: Record<string, any>; /* eslint-disable-line */
    start: () => void;
    pauseResume: () => void;
    reset: () => void;
    update: (v: number) => void;
  }
}
