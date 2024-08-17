import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest, { params }: { params: { id: string } }) {
  try {
    const { id } = params;

    const response = await fetch(`${process.env.SUBJECTS_API_BASE_URL}/api/subjects/user/${id}`);

    if (response.status === 204) {
      return NextResponse.json([], { status: 200 });
    }

    if (!response.ok) {
      throw new Error('Failed to fetch subjects for the user');
    }

    const data = await response.json();
    return NextResponse.json(data, { status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}