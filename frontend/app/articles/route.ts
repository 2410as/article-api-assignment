export async function GET() {
  const res = await fetch("http://localhost:8080/articles");
  const data = await res.json();
  return new Response(JSON.stringify(data), { status: 200 });
}

export async function POST() {
  const res = await fetch("http://localhost:8080/articles/import", {
    method: "POST",
  });
  return new Response(null, { status: res.status });
}
