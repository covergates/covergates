import { ReportState } from '.';

export function setCurrent(state: ReportState, report: Report): void {
  state.current = report;
}

export function startLoading(state: ReportState): void {
  state.loading = true;
}

export function stopLoading(state: ReportState): void {
  state.loading = false;
}
