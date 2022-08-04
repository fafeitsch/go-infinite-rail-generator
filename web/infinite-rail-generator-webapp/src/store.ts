import {BehaviorSubject, map, Observable, Subject, takeUntil} from 'rxjs';

interface State {
  seed: string
}

const state: State = {
  seed: ''
}

const state$ = new BehaviorSubject<State>(state)
const destroy$ = new Subject<void>()

export default {
  effect(observable$: Observable<any>){
    observable$.pipe(takeUntil(destroy$)).subscribe()
  },
  get: {
    seed$: state$.pipe(map(state => state.seed))
  },
  set: {
    seed(seed: string) {
      state$.next({...state$.value, seed: seed})
    }
  }
}
