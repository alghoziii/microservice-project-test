import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 50,          // jumlah virtual user
    duration: '10s',  // durasi tes
};

export default function () {
    const url = 'http://localhost:3001/api/orders';
    const payload = JSON.stringify({
        productId: 1,
        quantity: 1,
    });
    const headers = { 'Content-Type': 'application/json' };

    const res = http.post(url, payload, { headers });

    check(res, {
        'status 201': (r) => r.status === 201,
        'response time < 300ms': (r) => r.timings.duration < 300,
    });

    sleep(0.5);
}
