import { NextRequest, NextResponse } from 'next/server';

export async function POST(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const scanner = searchParams.get("scanner");

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_AGGREGATOR_API_URL}/run?scanner=${scanner?.toLowerCase()}`, {
      method: "POST"
    });
    if (!response.ok) {
      return NextResponse.json(
        { error: `Fetch error: ${response.status} ${response.statusText}` },
        { status: response.status }
      );
    }

    return NextResponse.json({ status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}