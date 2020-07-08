import { Server, Registry, Model, Factory, Response } from 'miragejs';

const User = Model.extend({
  login: '',
  email: ''
});

const Repository = Model.extend({
  URL: '',
  ReportID: '',
  NameSpace: '',
  Name: '',
  Branch: '',
  SCM: ''
});

const RepositoryFactory = Factory.extend({
  URL(i: number) {
    return `http://github/org${i}/repo${i}`;
  },
  ReportID(i: number) {
    return `report${i}`;
  },
  NameSpace(i: number) {
    return `org${i}`;
  },
  Name(i: number) {
    return `repo${i}`;
  },
  Branch: 'master',
  SCM: 'github'
});

const models = {
  user: User,
  repository: Repository
};

const factories = {
  repository: RepositoryFactory
};

type Models = typeof models;

function seeds(server: Server<Registry<Models, {}>>): void {
  server.schema.create('user', {
    login: 'blueworrybear',
    email: 'blueworrybear@gmail.com'
  });
  server.createList('repository', 5);
}

function routes(this: Server<Registry<Models, {}>>): void {
  this.namespace = '/api/v1';
  // user
  this.get('/user', schema => {
    const user = schema.first('user');
    return user !== null ? user.attrs : {};
  });
  // repository
  this.get('/repos/:scm', schema => {
    return schema.all('repository').models;
  });
  this.get('/repos/:scm/:namespace/:name/files', () => {
    const files = [];
    for (let i = 0; i < 10; i++) {
      files.push(`file${i}`);
    }
    files.push('dir/file.pl');
    files.push('main.pl');
    return files;
  });

  this.get('/repos/:scm/:namesapce/:name/content/*path', () => {
    return new Response(200, undefined, `print "hello";
my $s = "test";
if ($s =~ /^t/) {
\tprint 'match\\n';
}
    `);
  });

  this.get('/repos/:scm/:namespace/:name', (schema, request) => {
    const repo = schema.findBy('repository', { Name: request.params.name });
    return repo !== null ? repo.attrs : {};
  });
  // report
  this.get('/reports/:id', (_, request) => {
    const report: Report = {
      commit: '123456',
      reportID: `report${request.params.id}`,
      coverage: {
        Files: [
          {
            Name: 'main.pl',
            StatementCoverage: 0.8,
            StatementHits: [
              {
                LineNumber: 1,
                Hits: 1
              },
              {
                LineNumber: 2,
                Hits: 1
              }
            ]
          }
        ],
        StatementCoverage: 0.8
      }
    };
    return report;
  });
}

export class MockServer extends Server<Registry<Models, {}>> {
  constructor(environment = 'development') {
    super({
      environment,
      models: models,
      factories: factories,
      seeds: seeds,
      routes: routes
    });
  }
}

export function makeServer(environment = 'development'): MockServer {
  const server = new MockServer(environment);
  return server;
}
