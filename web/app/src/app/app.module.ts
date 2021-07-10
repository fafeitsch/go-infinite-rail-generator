import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {TrackTileViewerModule} from './track-tile-viewer/track-tile-viewer.module';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    TrackTileViewerModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
