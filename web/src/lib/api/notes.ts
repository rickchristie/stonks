import { apiPost } from './client';
import type { ApiResponse } from '$lib/types/api';

export type NoteStatus = 'Active' | 'Archived';

export type Note = {
	id: number;
	title: string;
	body: string;
	status: NoteStatus;
	createdTs: string;
	lastUpdatedTs: string;
};

export type CreateNoteResp = ApiResponse & {
	note: Note | null;
};

export type ListNotesResp = ApiResponse & {
	notes: Note[];
};

export type ArchiveNoteResp = ApiResponse & {
	note: Note | null;
};

export function createNote(title: string, body: string): Promise<CreateNoteResp> {
	return apiPost<CreateNoteResp>('/api/note/create', { title, body });
}

export function listNotes(): Promise<ListNotesResp> {
	return apiPost<ListNotesResp>('/api/note/list', {});
}

export function archiveNote(noteId: number): Promise<ArchiveNoteResp> {
	return apiPost<ArchiveNoteResp>('/api/note/archive', { noteId });
}
