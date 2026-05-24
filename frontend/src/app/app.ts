import { JsonPipe } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { firstValueFrom } from 'rxjs';

interface EmailAnalysis {
  id: string;
  email: string;
  domain: string;
  status: string;
  created_at: string;
  updated_at: string;
  timings?: {
    total_ms?: number;
    dns_ms?: number;
    whois_ms?: number;
    smtp_ms?: number;
    catch_all_ms?: number;
    disposable_ms?: number;
    typo_ms?: number;
    related_ms?: number;
    risk_ms?: number;
  };
  [key: string]: unknown;
}

@Component({
  selector: 'app-root',
  imports: [FormsModule, JsonPipe],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  private readonly http = inject(HttpClient);
  private timerId: ReturnType<typeof setInterval> | null = null;

  protected email = '';
  protected loading = false;
  protected error = '';
  protected analysis: EmailAnalysis | null = null;
  protected hasStarted = false;
  protected elapsedMs = 0;
  protected lastDurationMs = 0;

  protected get canSubmit(): boolean {
    return this.email.trim().length > 0 && !this.loading;
  }

  protected get loadingTimeText(): string {
    return `${(this.elapsedMs / 1000).toFixed(1)}s`;
  }

  protected get lastDurationText(): string {
    return `${(this.lastDurationMs / 1000).toFixed(2)}s`;
  }

  protected get loadingProgress(): number {
    // Simulated progress that never reaches 100% before the API response arrives.
    return Math.min(96, Math.round(14 + this.elapsedMs / 55));
  }

  protected get loadingStage(): string {
    if (this.elapsedMs < 700) {
      return 'Igniting scan core';
    }

    if (this.elapsedMs < 1600) {
      return 'Reading domain signals';
    }

    if (this.elapsedMs < 2800) {
      return 'Checking delivery paths';
    }

    return 'Waiting for API confirmation';
  }

  protected async analyze(): Promise<void> {
    const email = this.email.trim();

    if (!email) {
      this.error = 'Enter an email address first.';
      return;
    }

    this.hasStarted = true;
    this.loading = true;
    this.error = '';
    this.elapsedMs = 0;
    this.startTimer();
    const requestStart = performance.now();

    try {
      this.analysis = await firstValueFrom(
        this.http.post<EmailAnalysis>('/analyze', { email })
      );
      this.lastDurationMs = Math.round(performance.now() - requestStart);
    } catch {
      this.analysis = null;
      this.error = 'Could not reach the Go API at http://localhost:8080.';
      this.lastDurationMs = Math.round(performance.now() - requestStart);
    } finally {
      this.loading = false;
      this.stopTimer();
    }
  }

  protected newSearch(): void {
    this.stopTimer();
    this.hasStarted = false;
    this.analysis = null;
    this.error = '';
    this.loading = false;
    this.elapsedMs = 0;
  }

  private startTimer(): void {
    this.stopTimer();
    const startedAt = performance.now();
    this.timerId = setInterval(() => {
      this.elapsedMs = Math.round(performance.now() - startedAt);
    }, 80);
  }

  private stopTimer(): void {
    if (this.timerId) {
      clearInterval(this.timerId);
      this.timerId = null;
    }
  }
}
