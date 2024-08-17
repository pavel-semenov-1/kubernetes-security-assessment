import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest, { params }: { params: { id: string } }) {
  try {
    const { id } = params;

    const response = await fetch(`${process.env.USERS_API_BASE_URL}/api/users/${id}`);
    if (!response.ok) {
      throw new Error('Failed to fetch the user');
    }

    const data = await response.json();
    return NextResponse.json(data, { status: 200 });
  } catch (error: any) {
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}