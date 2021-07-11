import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TrackTileViewerComponent } from './track-tile-viewer.component';
import { DragDropModule } from '@angular/cdk/drag-drop';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {ReactiveFormsModule} from '@angular/forms';

@NgModule({
  declarations: [TrackTileViewerComponent],
  exports: [TrackTileViewerComponent],
  imports: [CommonModule, DragDropModule, LeafletModule, MatFormFieldModule, MatInputModule, ReactiveFormsModule],
})
export class TrackTileViewerModule {}
