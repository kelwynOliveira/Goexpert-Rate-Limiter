// @ts-nocheck
import http from 'k6/http';

export const options = {
  vus: 5,
  duration: '1m',
  thresholds: {
    'http_req_duration{status:200}': ['max>=0'],
    'http_req_duration{status:429}': ['max>=0'],
    'http_req_duration{status:500}': ['max>=0'],
  },
  'summaryTrendStats': ['min', 'med', 'avg', 'p(90)', 'p(95)', 'max', 'count'],
};

export default function () {
  http.get('http://rate-limiter_api:8080');

  http.get('http://rate-limiter_api:8080', {
    headers: {
      'API_KEY': 'any-api-key'
    }
  });
}
