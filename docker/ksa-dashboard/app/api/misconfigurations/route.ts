import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const reportId = searchParams.get("reportId");
  const search = searchParams.get("search");
  const resolved = searchParams.get("resolved");

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_PARSER_API_URL}/misconfigurations?reportId=${reportId}&search=${search || ""}&resolved=${resolved}`, {
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

export async function DELETE(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const reportId = searchParams.get("reportId");
  const search = searchParams.get("search");
  const id = searchParams.get("id");

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_PARSER_API_URL}/misconfigurations?reportId=${reportId}&search=${search || ""}&id=${id || ""}`, {
      method: "DELETE",
    });

    if (!response.ok) {
      return NextResponse.json({ error: `Failed to delete: ${response.statusText}` }, { status: response.status });
    }

    return NextResponse.json({ message: "Delete successful" }, { status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message || "Internal server error" }, { status: 500 });
  }
}

export async function PATCH(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const reportId = searchParams.get("reportId");
  const search = searchParams.get("search");
  const id = searchParams.get("id");

  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_PARSER_API_URL}/misconfigurations?reportId=${reportId}&search=${search || ""}&id=${id || ""}`, {
      method: "PATCH",
    });

    if (!response.ok) {
      return NextResponse.json({ error: `Failed to patch: ${response.statusText}` }, { status: response.status });
    }

    return NextResponse.json({ message: "Patch successful" }, { status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message || "Internal server error" }, { status: 500 });
  }
}