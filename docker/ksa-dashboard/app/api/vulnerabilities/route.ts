import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const severity = searchParams.get("severity");
  const scanner = searchParams.get("scanner");

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_PARSER_API_URL}/vulnerabilities?scanner=${scanner?.toLowerCase()}&severity=${severity}&reportId=1`, {
      method: "GET"
    });
    if (!response.ok) {
      return NextResponse.json(
        { error: `Fetch error: ${response.status} ${response.statusText}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    const mappedData = data.map((item: any) => ({
      ...item,
      ID: item.VulnerabilityID
    }))
    return NextResponse.json(mappedData, { status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}