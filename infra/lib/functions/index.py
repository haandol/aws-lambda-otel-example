import json
import requests
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.sampling import ALWAYS_ON

# set the sampler onto the global tracer provider
trace.set_tracer_provider(TracerProvider(sampler=ALWAYS_ON))

def handler(event, context):
    print(json.dumps(event))

    tracer = trace.get_tracer(__name__)
    with tracer.start_span("requests amazon") as span:
        url = "https://aws.amazon.com/"
        requests.get(url, { 'timeout': 1, })
        span.set_attribute("http.method", "GET")
        span.set_attribute("http.url", url)

    with tracer.start_span("devide zero") as span:
        try:
            1 / 0
        except ZeroDivisionError as err:
            span.record_exception(err)
            raise err

    return "ok"
