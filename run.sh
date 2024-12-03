#!/bin/bash

# Arguments
export K6_WEB_DASHBOARD=true
export K6_WEB_DASHBOARD_PORT=-1
export K6_WEB_DASHBOARD_EXPORT=html-report.html

# Execute
xk6 run --out json=result.json ./examples/publish-listen.js
