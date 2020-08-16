declare var VUE_BASE: any; // eslint-disable-line

declare interface User {
  login?: string;
  email?: string;
  avatar?: string;
}

declare interface Repository {
  URL: string;
  SCM: string;
  ReportID?: string;
  NameSpace: string;
  Name: string;
  Branch: string;
}

declare interface RepositorySetting {
  filters?: string[];
  mergePR?: boolean;
}

declare interface Commit {
  sha: string;
  committer: string;
  committerAvatar: string;
  message: string;
}

declare interface StatementHit {
  LineNumber: number;
  Hits: number;
}

declare interface SourceFile {
  Name: string;
  StatementCoverage: number;
  StatementHits: StatementHit[];
}

declare interface Coverage {
  files?: SourceFile[];
  type: string;
  statementCoverage: number | 0;
}

declare interface Report {
  ref?: string;
  commit: string;
  coverages: Coverage[];
  reportID: string;
  files?: string[];
  createdAt?: string;
}
