import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { ToolbarComponent } from './toolbar/toolbar.component';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import { HttpClientModule } from '@angular/common/http';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, ToolbarComponent, MatSlideToggleModule, HttpClientModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'commanderer-frontend';
  enabled = false;

  constructor(private http: HttpClient) {}

  onToggleChange(event: any) {
    this.enabled = event.checked;
    const newMode = { enabled: this.enabled };

    this.http.put('/api/mode', newMode)
      .subscribe(
        response => {
          console.log('Toggle status updated:', response);
        },
        error => {
          console.error('Error updating toggle status:', error);
        }
      );
  }
}
