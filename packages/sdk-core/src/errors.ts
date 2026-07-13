export type ProblemDetails = {
  type?: string;
  title?: string;
  status?: number;
  detail?: string;
  instance?: string;
};

export class FinOpsHttpError extends Error {
  readonly status: number;
  readonly problem?: ProblemDetails;

  constructor(message: string, status: number, problem?: ProblemDetails) {
    super(message);
    this.name = "FinOpsHttpError";
    this.status = status;
    this.problem = problem;
  }

  static async fromResponse(response: Response): Promise<FinOpsHttpError> {
    let problem: ProblemDetails | undefined;

    try {
      const contentType = response.headers.get("content-type") ?? "";
      if (
        contentType.includes("application/json") ||
        contentType.includes("application/problem+json")
      ) {
        problem = (await response.json()) as ProblemDetails;
      }
    } catch {
      // Ignore malformed error bodies.
    }

    const detail = problem?.detail ?? problem?.title ?? response.statusText;
    return new FinOpsHttpError(
      detail || `Request failed with status ${response.status}`,
      response.status,
      problem,
    );
  }
}
