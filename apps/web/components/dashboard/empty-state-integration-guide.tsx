import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export function EmptyStateIntegrationGuide() {
  return (
    <Card className="border-dashed">
      <CardHeader>
        <CardTitle>No usage events yet</CardTitle>
        <CardDescription>
          Connect your application with the AI FinOps SDK to start tracking LLM spend.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <pre className="overflow-x-auto rounded-lg bg-muted p-4 text-sm">
          {`import { FinOpsClient } from '@ai-finops/sdk';

const client = new FinOpsClient({
  apiKey: process.env.AI_FINOPS_API_KEY,
  baseUrl: '${process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8888/v1"}',
});

await client.ingest({
  provider: 'openai',
  model: 'gpt-4o',
  input_tokens: 1200,
  output_tokens: 340,
  occurred_at: new Date().toISOString(),
});`}
        </pre>
      </CardContent>
    </Card>
  );
}
