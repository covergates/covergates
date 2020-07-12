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

export function setSource(state: ReportState, source: string): void {
  state.source = source;
}

export function setError(state: ReportState, error: Error): void {
  state.error = error;
}
