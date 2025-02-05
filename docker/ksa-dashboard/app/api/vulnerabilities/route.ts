import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const scanner = searchParams.get("scanner");
  const reportId = searchParams.get("reportId");
  const search = searchParams.get("search");

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_PARSER_API_URL}/vulnerabilities?scanner=${scanner?.toLowerCase()}&reportId=${reportId}&search=${search || ""}`, {
      method: "GET"
    });
    if (!response.ok) {
      return NextResponse.json(
        { error: `Fetch error: ${response.status} ${response.statusText}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data, { status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}