import { NextResponse } from 'next/server';

export async function GET() {
    return NextResponse.json({ websocketUrl: process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8124/ws" }, { status: 200 });
}